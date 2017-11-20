#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 MVP;
uniform mat4 model;

out vec2 TexCoord;
out vec3 FragPos;
out vec3 Normal;

void main() {
    gl_Position = MVP * vec4(position, 1.0);
    TexCoord = aTexCoord;
    FragPos = vec3(model * vec4(position, 1.0));

    Normal = mat3(transpose(inverse(model))) * normal;
}
