#version 410 core

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aTangent;

uniform mat4 MVP;
uniform mat4 MV;
uniform mat4 view;
uniform mat4 InvView;
uniform mat4 model;

#include "pbr_lights.glsl"
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
    mat3 TBN;
    vec4 FragPosLightSpace;
} vs_out;

void main() {

    // the position of the fragment in the perspective space
    gl_Position = MVP * vec4(aPosition, 1.0);

    // the position of the camara relative to the fragment
    vs_out.V_Pos = vec3(MV * vec4(aPosition, 1.0));

    vs_out.TexCoord = aTexCoord ;

    //surface normal in the world space, used for lookup env map coordinates
    vs_out.Normal = mat3(model) * aNormal;

    vec3 eyeDir = normalize(vec3(InvView * vec4(normalize(-vs_out.V_Pos), 0.0)));
    vs_out.Reflection = reflect(-eyeDir, vs_out.Normal);

    vec3 N = normalize(vec3(MV * vec4(aNormal, 0.0)));
    vec3 T = normalize(vec3(MV * vec4(aTangent, 0.0)));
    vs_out.TBN = mat3(T, cross(N, T), N);

    vs_out.FragPosLightSpace = LightVP * vec4(vec3(model * vec4(aPosition, 1.0)), 1.0);

    // transform light positions into view space
    for (int i = 0; i < numLights; i++ ) {
        vs_out.V_LightPositions[i] = vec3(view * vec4(lights[i].position, 1));
    }
}
