#include "include.glsl"

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out VS_OUT
{
    out vec2 TexCoord;
    out vec3 Normal;
} vs_out;

void setOutputs() {
    gl_Position = projection * view * model * vec4(aPosition, 1.0);
    vs_out.Normal = aNormal;
    vs_out.TexCoord = aTexCoord;
}
