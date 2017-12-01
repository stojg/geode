#version 410 core

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 MVP;
uniform mat4 MV;
uniform mat4 InverseMV;
uniform mat4 view;


struct Light {
    vec3 position;
    vec3 color;
};

uniform Light pointLights[16];
uniform int numPointLights;

out VS_OUT
{
    vec3 V_Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[16];
    vec3 W_ViewPos;
} vs_out;

void main() {

    // the position of the fragment in the perspective space
    gl_Position = MVP * vec4(aPosition, 1.0);

    // the position of the the view relative to the fragment
    vs_out.W_ViewPos = vec3(MV * vec4(aPosition, 1.0));

    vs_out.TexCoord = aTexCoord;

    // transform normals into view space
    vs_out.V_Normal = normalize(mat3(InverseMV) * aNormal);

    // transform light positions into view space
    for (int i = 0; i < numPointLights; i++ ) {
        // point lights have a position, so it's vec4(pos, 1); directional lights are vec4(pos, 0);
        vs_out.V_LightPositions[i] = vec3(view * vec4(pointLights[i].position, 1));
    }
}
