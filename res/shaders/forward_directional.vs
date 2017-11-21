#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

// shadow
uniform mat4 lightSpaceMatrix;

uniform vec3 lightPos;

out vec2 TexCoord;
out vec3 FragPos;
out vec3 Normal;
out vec3 LightPos;

// shadow
out vec4 FragPosLightSpace;

void main() {
    gl_Position = projection * view * model * vec4(position, 1.0);
    FragPos = vec3(view  * model * vec4(position, 1.0));
    Normal = mat3(transpose(inverse(view * model))) * normal;
    LightPos = vec3(view * vec4(lightPos, 0.0));
    TexCoord = aTexCoord;

    // shadow
    FragPosLightSpace = lightSpaceMatrix * model * vec4(position, 1.0);
}
