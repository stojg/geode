#version 410 core

in VS_OUT
{
    vec3 V_Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[16];
    vec3 W_ViewPos;
} vs_in;


uniform float specularStrength = 0.1;

struct Material {
    vec3 albedo;
    float metallic;
    float roughness;
};

uniform Material material;

#include "point_lights.glsl"
uniform int numPointLights;
uniform Light pointLights[16];

out vec4 FragColor;

vec3 CalcPointLight(vec3 lightPosition, Light light, Material material, vec3 norm, vec3 viewPos) {
    vec3 lightDiff = lightPosition - viewPos;
    float distance = length(lightDiff);

    if (distance > light.distance) {
        return vec3(0);
    }
    vec3 lightDirection = normalize(lightDiff);

    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    // diffuse
    float diff = max(dot(norm, lightDirection), 0.0);

    // specular
    vec3 halfwayDir = normalize(lightDirection - normalize(viewPos));
    vec3 reflectDir = reflect(-lightDirection, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 16);

    // combine results
    vec3 diffuseColor = light.color * diff * material.albedo;
    vec3 specularColor = light.color * spec * material.albedo;

    diffuseColor *= attenuation;
    specularColor *= attenuation;

    return diffuseColor + specularColor;
}

void main() {
    vec3 normal = normalize(vs_in.V_Normal);

    float ambientStrength = 0.01;
    vec3 final = ambientStrength * material.albedo;

    for (int i = 0; i < numPointLights; i++) {
        final += CalcPointLight(vs_in.V_LightPositions[i], pointLights[i], material, normal, vs_in.W_ViewPos);
    }

    FragColor = vec4(final, 0);
}
