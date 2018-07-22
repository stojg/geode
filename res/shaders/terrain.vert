#version 410 core

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aTangent;

uniform mat4 model;

#include "matrices.glsl"
#include "light_struct.glsl"
uniform Light lights[16];
uniform int numLights;
uniform mat4 LightVP;

out VS_OUT
{
    vec3 Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[8];
    vec3 V_Pos;
    vec3 Reflection;
    vec3 W_Pos;
    mat3 TBN;
    mat3 W_TBN;
    vec4 FragPosLightSpace;
} vs_out;

void main() {

    mat4 MV = view * model;
    mat4 MVP = projection * MV;

    // the position of the fragment in the perspective space
    gl_Position = MVP * vec4(aPosition, 1.0);

    // the position of the camara relative to the fragment
    vs_out.V_Pos = vec3(MV * vec4(aPosition, 1.0));

    vs_out.TexCoord = aTexCoord * 128;

    // surface normal in the world space, used for lookup env map coordinates
    vs_out.Normal = mat3(model) * aNormal;

    vec3 eyeDir = normalize(vec3(InvView * vec4(normalize(-vs_out.V_Pos), 0.0)));
    vs_out.Reflection = reflect(-eyeDir, vs_out.Normal);

    vec4 pos = model * vec4(aPosition,1);
    vs_out.W_Pos = vec3(pos.xyz) / pos.w;

    vec3 N = mat3(MV) * (aNormal);
    vec3 T = mat3(MV) * (aTangent);
    // re-orthogonalize T with respect to N
    T = normalize(T - dot(T, N) * N);
    vs_out.TBN = mat3(T, cross(N, T), N);

    vec3 W_N = mat3(model) * (aNormal);
    vec3 W_T = mat3(model) * (aTangent);
    // re-orthogonalize W_T with respect to W_N
    W_T = normalize(W_T - dot(W_T, W_N) * W_N);
    vs_out.W_TBN = mat3(W_T, cross(W_N, W_T), W_N);

    vs_out.FragPosLightSpace = LightVP * vec4(vec3(model * vec4(aPosition, 1.0)), 1.0);

    // transform light positions into view space
    for (int i = 0; i < numLights; i++ ) {
        vs_out.V_LightPositions[i] = vec3(view * vec4(lights[i].position, 1));
    }

    vs_out.FragPosLightSpace = LightVP * vec4(vec3(model * vec4(aPosition, 1.0)), 1.0);

    // transform light positions into view space
    for (int i = 0; i < numLights; i++ ) {
        vs_out.V_LightPositions[i] = vec3(view * vec4(lights[i].position, 1));
    }
}
