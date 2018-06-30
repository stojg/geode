#version 410 core

struct Material {
    sampler2D albedo;
    sampler2D metallic;
    sampler2D roughness;
    sampler2D normal;
};

in VS_OUT
{
    // surface normal in the world space
    vec3 Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[8];
    // camera position in world space
    vec3 V_Pos;
    vec3 Reflection;
    mat3 TBN;
    vec4 FragPosLightSpace;
} vs_in;

#include "light_struct.glsl"
#include "pbr.glsl"
#include "pbr_ambient.glsl"
#include "shadow.glsl"

void main() {

    vec3 normal = texture(material.normal, vs_in.TexCoord).rgb;
    // transform normal vector to range [-1,1]
    normal = normalize(normal * 2.0 - 1.0);
    normal = normalize(vs_in.TBN * normal);

    Mtrl mtrl;
    mtrl.albedo = texture(material.albedo, vs_in.TexCoord).rgb;
    mtrl.metallic = texture(material.metallic, vs_in.TexCoord).r;
    mtrl.roughness = texture(material.roughness, vs_in.TexCoord).r;

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, mtrl.albedo, mtrl.metallic);

    vec3 Lo = vec3(0.0);

    vec3 V = normalize(-vs_in.V_Pos);

    for (int i = 0; i < numLights; i++) {
        if (lights[i].constant == 0) {
            Lo += CalcDirectional(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
            Lo *= ShadowCalculation(vs_in.FragPosLightSpace);
        } else if (lights[i].cutoff > 0) {
            Lo += CalcSpot(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
        } else {
            Lo += CalcPoint(F0, vs_in.V_LightPositions[i], lights[i], mtrl, normal, vs_in.V_Pos, V);
        }
    }

    if (x_enable_env_map == 1) {
        Lo += CalcAmbient(normal, V, F0, mtrl);
    }

    FragColor = vec4(Lo, 1);
}
