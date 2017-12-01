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

#include "pbr.glsl"

vec3 CalcPointLight(vec3 F0, vec3 lightPosition, Light light, Material material, vec3 norm, vec3 viewPos) {
    vec3 lightDiff = lightPosition - viewPos;
    float distance = length(lightDiff);

    if (distance > light.distance) {
        return vec3(0);
    }

    vec3 lightDirection = normalize(lightDiff);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    vec3 V = normalize(-viewPos);
    vec3 N = norm;
    vec3 L = normalize(lightPosition - viewPos);
    vec3 H = normalize(V + L);

    vec3 radiance = light.color * attenuation;

    // cook-torrance brdf
    float NDF = DistributionGGX(N, H, material.roughness);
    float G = GeometrySmith(N, V, L, material.roughness);
    vec3 F = fresnelSchlick(max(dot(H, V), 0.0), F0);

    vec3 kS = F;
    vec3 kD = vec3(1.0) - kS;
    kD *= 1.0 - material.metallic;

    vec3 nominator    = NDF * G * F;
    float denominator = 4.0 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0);
    vec3 specular     = nominator / max(denominator, 0.001);

    // add to outgoing radiance Lo
    float NdotL = max(dot(N, L), 0.0);
    return (kD * material.albedo / PI + specular) * radiance * NdotL;
}

void main() {
    vec3 normal = normalize(vs_in.V_Normal);

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, material.albedo, material.metallic);

    // reflectance equation
    vec3 Lo = vec3(0.0);

    for (int i = 0; i < numPointLights; i++) {
        Lo += CalcPointLight(F0, vs_in.V_LightPositions[i], pointLights[i], material, normal, vs_in.W_ViewPos);
    }

    float ambientStrength = 0.01;
    vec3 color = ambientStrength * material.albedo + Lo;

    FragColor = vec4(color, 1);
}
