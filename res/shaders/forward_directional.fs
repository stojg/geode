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

uniform sampler2D diffuse;
uniform DirectionalLight directionalLight;

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 FragPos;

out vec4 fragColor;

float specularStrength = 0.5;

// shadow
in vec4 FragPosLightSpace;
uniform sampler2D x_shadowMap;

float sampleShadowMap(sampler2D shadowMap, vec2 coords, float compare, float bias)
{
    vec2 samplingCoords = coords;
    return step(texture(shadowMap, samplingCoords).r, compare - bias);
}

float sampleShadowMapLinear(sampler2D shadowMap, vec2 coords, float compare, float bias, vec2 texelSize)
{
    vec2 pixelPos = coords / texelSize + vec2(0.5);
    vec2 fractPart = fract(pixelPos);
    vec2 startTexel = (pixelPos - fractPart) * texelSize;

    float blTexel = sampleShadowMap(shadowMap, startTexel, compare, bias);
    float brTexel = sampleShadowMap(shadowMap, startTexel + vec2(texelSize.x, 0.0), compare, bias);
    float tlTexel = sampleShadowMap(shadowMap, startTexel + vec2(0.0, texelSize.y), compare, bias);
    float trTexel = sampleShadowMap(shadowMap, startTexel + texelSize, compare, bias);

    float mixA = mix(blTexel, tlTexel, fractPart.y);
    float mixB = mix(brTexel, trTexel, fractPart.y);
    return mix(mixA, mixB, fractPart.x);
}

float sampleShadowMapPCF(sampler2D shadowMap, vec2 coords, float compare, float bias, vec2 texelSize)
{
    const float NUM_SAMPLES = 3.0;
    const float SAMPLES_START = (NUM_SAMPLES-1.0)/2.0;
    const float NUM_SAMPLES_SQUARED = NUM_SAMPLES * NUM_SAMPLES;

    float shadow = 0.0;
    for(float x = -SAMPLES_START; x <= SAMPLES_START; x += 1.0) {
        for(float y = -SAMPLES_START; y <= SAMPLES_START; y += 1.0) {
            vec2 offset = vec2(x,y) * texelSize;
            shadow += sampleShadowMapLinear(shadowMap, coords + offset , compare, bias, texelSize);
        }
    }
    return shadow /= NUM_SAMPLES_SQUARED;
}

float ShadowCalculation(vec4 fragPosLightSpace, vec3 normal, vec3 lightDir)
{
    // perform perspective divide, since it's not done automatically for us
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
    // transform from [-0.5,0.5] to [0,1] range so we can use it for sampling
    projCoords = projCoords * 0.5 + 0.5;

    // dont shadow things outside the light frustrum far plane
    if(projCoords.z > 1.0) {
        return 0.0;
    }

    float bias = max(0.02 * (1.0 - dot(normal, lightDir)), 0.02);
    vec2 texelSize = 0.5 / textureSize(x_shadowMap, 0);

    //return sampleShadowMap(x_shadowMap, projCoords.xy, projCoords.z, bias);
    //return sampleShadowMapLinear(x_shadowMap, projCoords.xy, projCoords.z, bias, texelSize);
    return sampleShadowMapPCF(x_shadowMap, projCoords.xy, projCoords.z, bias, texelSize);
}

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDir = normalize(LightPos);

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * directionalLight.base.color;

    vec3 viewDir = normalize(-FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * directionalLight.base.color;

    // calculate shadow
    float shadow = ShadowCalculation(FragPosLightSpace, norm, lightDir);

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= (1-shadow);
}
