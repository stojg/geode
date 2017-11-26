#version 410 core

struct Attenuation
{
    float constant;
    float linear;
    float exponent;
};

struct BaseLight
{
    vec3 color;
};

struct DirectionalLight
{
    BaseLight base;
    vec3 direction;
};

struct PointLight
{
    BaseLight base;
    Attenuation atten;
    vec3 position;
};

struct SpotLight
{
    PointLight pointLight;
    vec3 direction;
    float cutoff;
};

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 ModelViewPos;

out vec4 fragColor;

const float specularStrength = 0.5;

uniform sampler2D diffuse;
uniform SpotLight spotLight;
uniform float x_varianceMin;
uniform float x_lightBleedReductionAmount;

uniform mat4 view;

// shadow
in vec4 FragPosLightSpace;
uniform sampler2D x_shadowMap;

float linstep(float low, float high, float v)
{
	return clamp((v-low)/(high-low), 0.0, 1.0);
}

float sampleVarianceShadowMap(sampler2D shadowMap, vec2 coords, float compare, float varianceMin, float lightBleedReductionAmount)
{
    // return step(compare, texture(shadowMap, coords.xy).r);
    vec2 moments = texture(shadowMap, coords).xy;
	float p = step(compare, moments.x);

	float variance = max(moments.y - moments.x * moments.x, varianceMin);

	float d = compare - moments.x;
	float pMax = linstep(lightBleedReductionAmount, 1.0, variance / (variance + d*d));

	return min(max(p, pMax), 1.0);
}

float ShadowCalculation(vec4 fragPosLightSpace, vec3 normal, vec3 lightDir, float varianceMin, float lightBleedReductionAmount)
{
    // perform perspective divide, since it's not done automatically for us
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
    // transform from [-0.5,0.5] to [0,1] range so we can use it for sampling
    projCoords = projCoords * 0.5 + 0.5;

    // dont shadow things outside the light frustrum far plane
    if(projCoords.z > 1.0) {
        return 1.0;
    }

    return sampleVarianceShadowMap(x_shadowMap, projCoords.xy, projCoords.z, varianceMin, lightBleedReductionAmount);
}

void main() {

    vec3 norm = Normal;
    vec3 color = spotLight.pointLight.base.color;

    vec3 lightDiff = LightPos - ModelViewPos;
    float lightDistance = length(lightDiff);
    vec3 lightDir = normalize(lightDiff);

    vec3 spotlightDirection = (view * vec4(spotLight.direction,0)).xyz;
    float theta = dot(lightDir, normalize(-spotlightDirection));
    if(theta < spotLight.cutoff) {
        fragColor = vec4(0,0,0,0);
        return;
    }

    float attenuation = 1.0 / (spotLight.pointLight.atten.constant + spotLight.pointLight.atten.linear * lightDistance + spotLight.pointLight.atten.exponent * (lightDistance * lightDistance));

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * color;

    vec3 halfwayDir = normalize(lightDir - ModelViewPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * color;

    float shadow = ShadowCalculation(FragPosLightSpace, norm, lightDir, x_varianceMin, x_lightBleedReductionAmount);

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= attenuation;
    fragColor *= shadow;
}
