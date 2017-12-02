#version 410 core

in VS_OUT
{
    vec3 Normal;
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

uniform samplerCube x_irradianceMap;

uniform Material material;

#include "point_lights.glsl"
uniform int numPointLights;
uniform Light pointLights[16];

out vec4 FragColor;

vec3 CalcPointLight(vec3 F0, vec3 lightPosition, Light light, Material material, vec3 norm, vec3 viewPos);

#include "pbr.glsl"

vec3 fresnelSchlickRoughness(float cosTheta, vec3 F0, float roughness)
{
    return F0 + (max(vec3(1.0 - roughness), F0) - F0) * pow(1.0 - cosTheta, 5.0);
}

void main() {
    vec3 normal = normalize(vs_in.V_Normal);

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, material.albedo, material.metallic);

    vec3 Lo = vec3(0.0);
    for (int i = 0; i < numPointLights; i++) {
        Lo += CalcPointLight(F0, vs_in.V_LightPositions[i], pointLights[i], material, normal, vs_in.W_ViewPos);
    }

    vec3 V = normalize(-vs_in.W_ViewPos);
    float nDotV = max(dot(vs_in.V_Normal, V), 0.0);
    vec3 kS = fresnelSchlickRoughness(nDotV, F0, material.roughness);
    vec3 kD = 1.0 - kS;

    vec3 irradiance = texture(x_irradianceMap, vs_in.Normal).rgb;
    vec3 diffuse    = irradiance * material.albedo;
    vec3 ambient    = (kD * diffuse);

    vec3 color = ambient + Lo;

    FragColor = vec4(color, 1);
}


