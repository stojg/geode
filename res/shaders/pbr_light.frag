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

vec3 CalcPointLight(vec3 F0, vec3 lightPosition, Light light, Material material, vec3 norm, vec3 viewPos);

void main() {
    vec3 normal = normalize(vs_in.V_Normal);

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, material.albedo, material.metallic);

    vec3 Lo = vec3(0.0);
    for (int i = 0; i < numPointLights; i++) {
        Lo += CalcPointLight(F0, vs_in.V_LightPositions[i], pointLights[i], material, normal, vs_in.W_ViewPos);
    }

    float ambientStrength = 0.02;
    vec3 color = ambientStrength * material.albedo + Lo;

    FragColor = vec4(color, 1);
}

#include "pbr.glsl"
