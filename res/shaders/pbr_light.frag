#version 410 core

struct Material {
    sampler2D albedo;
    sampler2D metallic;
    sampler2D roughness;
};

uniform sampler2D   x_brdfLUT;
uniform samplerCube x_irradianceMap;
uniform samplerCube x_prefilterMap;

uniform mat4 view;

in VS_OUT
{
    // surface normal in the world space
    vec3 Normal;
    // surface normal in view space
    vec3 V_Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[8];
    // camera position in world space
    vec3 V_Pos;
    vec3 Reflection;
} vs_in;

#include "pbr_lights.glsl"
#include "pbr.glsl"

void main() {

    // Normal in view space
    vec3 normal = normalize(vs_in.V_Normal);

    Mtrl mtrl;

    mtrl.albedo = texture(material.albedo, vs_in.TexCoord).rgb;
    mtrl.metallic = texture(material.metallic, vs_in.TexCoord).r;
    mtrl.roughness = texture(material.roughness, vs_in.TexCoord).r;


    vec3 albedo = mtrl.albedo;
    vec3 F0 = vec3(0.04);
    F0 = mix(F0, albedo, mtrl.metallic);

    vec3 Lo = vec3(0.0);

    vec3 V = normalize(-vs_in.V_Pos);

    for (int i = 0; i < numLights; i++) {
        if (lights[i].constant == 0) {
            Lo += CalcDirectional(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
        } else if (lights[i].cutoff > 0) {
            Lo += CalcSpot(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
        } else {
            Lo += CalcPoint(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
        }
    }

    if (x_enable_env_map == 0) {
        FragColor = vec4(Lo, 1);
        return;
    }

    vec3 F = fresnelSchlickRoughness(max(dot(normal, V), 0.0), F0, mtrl.roughness);

    vec3 kS = F;
    vec3 kD = 1.0 - kS;
    kD *= 1.0 - mtrl.metallic;

    // diffuse
    vec3 irradiance = texture(x_irradianceMap, vs_in.Normal).rgb;
    vec3 diffuse    = irradiance * albedo;

    // specular
    const float MAX_REFLECTION_LOD = 4.0;
    vec3 prefilteredColor = textureLod(x_prefilterMap, vs_in.Reflection,  mtrl.roughness * MAX_REFLECTION_LOD).rgb;
    vec2 brdf  = texture(x_brdfLUT, vec2(max(dot(normal, V), 0.0), mtrl.roughness)).rg;
    vec3 specular = prefilteredColor * (F * brdf.x + brdf.y);

    // sum up all ambient
    vec3 ambient = (kD * diffuse + specular);

    // combine with lights
    vec3 color = Lo + ambient;

    FragColor = vec4(color, 1);
}
