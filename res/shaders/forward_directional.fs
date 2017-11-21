#version 410 core

uniform sampler2D diffuse;
uniform vec3 lightColor;

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 FragPos;

out vec4 fragColor;

float specularStrength = 0.5;

// shadow
in vec4 FragPosLightSpace;
uniform sampler2D shadowMap;
uniform mat4 lightSpaceMatrix;
uniform mat4 view;

float ShadowCalculation(vec4 fragPosLightSpace, vec3 normal)
{
    // perform perspective divide
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
    // transform to [0,1] range
    projCoords = projCoords * 0.5 + 0.5;

    // dont shadow things outside the light frustrum far plane
    if(projCoords.z > 1.0) {
        return 0.0;
    }

    // @todo pass this in?
    vec3 lightDir = normalize(LightPos);

    float shadow = 0.0;
    float bias = max(0.05 * (1.0 - dot(normal, lightDir)), 0.001);

    vec2 texelSize = 0.5 / textureSize(shadowMap, 0);
    // Percentage Closing Filter
    for(int x = -1; x <= 1; ++x) {
        for(int y = -1; y <= 1; ++y) {
            float pcfDepth = texture(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r;
            shadow += projCoords.z - bias > pcfDepth ? 1.0 : 0.0;
        }
    }
    return shadow /= 9.0;
}

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDir = normalize(LightPos);

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * lightColor;

    vec3 viewDir = normalize(-FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * lightColor;

    // calculate shadow
    float shadow = ShadowCalculation(FragPosLightSpace, norm);

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= (1-shadow);
}
