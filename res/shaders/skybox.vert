#version 410 core

layout (location = 0) in vec3 aPos;

out vec3 TexCoords;

#include "matrices.glsl"

void main()
{
    TexCoords = aPos;
    // mat4(mat3(view)) removes rotation from the view
    vec4 pos = projection * mat4(mat3(view)) * vec4(aPos, 1.0);
    gl_Position = pos.xyww;
}
