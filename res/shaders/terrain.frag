#version 410 core

uniform sampler2D albedo;
uniform sampler2D metallic;
uniform sampler2D roughness;
uniform sampler2D normal;

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

#include "matrices.glsl"
#include "light_struct.glsl"
#include "pbr.glsl"
#include "pbr_ambient.glsl"
#include "shadow.glsl"
#include "fog.glsl"

void main() {

    vec3 N = texture(normal, vs_in.TexCoord).rgb;
    // transform normal vector to range [-1,1]
    N = normalize(N * 2.0 - 1.0);
    N = normalize(vs_in.TBN * N);

    Mtrl mtrl;
    mtrl.albedo = texture(albedo, vs_in.TexCoord).rgb;
    mtrl.metallic = texture(metallic, vs_in.TexCoord).r;
    mtrl.roughness = texture(roughness, vs_in.TexCoord).r;

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, mtrl.albedo, mtrl.metallic);

    vec3 Lo = vec3(0.0);

    vec3 V = normalize(-vs_in.V_Pos);

    for (int i = 0; i < numLights; i++) {
        if (lights[i].constant == 0) {
            Lo += CalcDirectional(F0, vs_in.V_LightPositions[i], lights[i], mtrl, N, vs_in.V_Pos, V);
            Lo *= ShadowCalculation(vs_in.FragPosLightSpace);
        } else if (lights[i].cutoff > 0) {
            Lo += CalcSpot(F0, vs_in.V_LightPositions[i], lights[i], mtrl, N, vs_in.V_Pos, V);
        } else {
            Lo += CalcPoint(F0, vs_in.V_LightPositions[i], lights[i], mtrl, N, vs_in.V_Pos, V);
        }
    }

    if (x_enable_env_map == 1) {
        Lo += CalcAmbient(N, V, F0, mtrl);
    }

    Lo = fogCalc(Lo, vs_in.V_Pos);

    FragColor = vec4(Lo, 1);
}
