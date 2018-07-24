#version 410 core

uniform sampler2D albedo;
uniform sampler2D albedo2;
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
    vec3 W_Pos;
    vec3 Reflection;
    mat3 TBN;
    mat3 W_TBN;
    vec4 FragPosLightSpace;
} vs_in;

#include "matrices.glsl"
#include "light_struct.glsl"
#include "pbr.glsl"
#include "pbr_ambient.glsl"
#include "shadow.glsl"
#include "fog.glsl"

uniform vec3 x_camPos;

void main() {

    vec3 N = texture(normal, vs_in.TexCoord).rgb;
    // transform normal vector to range [-1,1]
    N = normalize(N * 2.0 - 1.0);
    vec3 V_N = normalize(vs_in.TBN * N);
    vec3 W_N = normalize(vs_in.W_TBN * N);

    Mtrl mtrl;
    float factor = acos(dot(vec3(0,1,0), normalize(vs_in.Normal)))*2/3.1415;
    factor = clamp(factor, 0.0, 1.0);
    mtrl.albedo = mix(texture(albedo, vs_in.TexCoord).rgb, texture(albedo2, vs_in.TexCoord).rgb, factor);
    mtrl.metallic = texture(metallic, vs_in.TexCoord).r;
    mtrl.roughness = texture(roughness, vs_in.TexCoord).r;

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, mtrl.albedo, mtrl.metallic);

    vec3 Lo = vec3(0.0);

    vec3 V = normalize(-vs_in.V_Pos);
    for (int i = 0; i < numLights; i++) {
        if (lights[i].constant == 0) {
            Lo += CalcDirectional(F0, vs_in.V_LightPositions[i], lights[i], mtrl, V_N, vs_in.V_Pos, V);
            float shadow = ShadowCalculation(vs_in.FragPosLightSpace);
            vec3 cshadow = pow( vec3(shadow), vec3(1.0, 1.2, 1.5) );
            Lo *= cshadow;
        } else if (lights[i].cutoff > 0) {
            Lo += CalcSpot(F0, vs_in.V_LightPositions[i], lights[i], mtrl, V_N, vs_in.V_Pos, V);
        } else {
            Lo += CalcPoint(F0, vs_in.V_LightPositions[i], lights[i], mtrl, V_N, vs_in.V_Pos, V);
        }
    }

    vec3 vv = normalize(x_camPos - vs_in.W_Pos);
    vec3 Reflection = -reflect(vv, W_N);
    if (x_enable_env_map == 1) {
        Lo += CalcAmbient(W_N, vv, F0, mtrl, Reflection) * 0.2;
    }

    Lo = fogCalc(Lo, vs_in.V_Pos);

    FragColor = vec4(Lo, 1);
}
