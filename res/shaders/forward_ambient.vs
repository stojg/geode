#version 410

layout (location = 0) in vec3 position;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 MVP;

out vec2 TexCoord;

void main() {
    gl_Position = MVP * vec4(position, 1.0);
    TexCoord = aTexCoord;
}
