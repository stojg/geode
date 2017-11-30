#version 410 core

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform vec3 x_lightPositions[16];
uniform int x_numPointLights;

out vec3 Normal;
out vec2 TexCoord;
out vec3 LightPositions[16];
out vec3 FragPos;
out vec3 ModelViewPos;

void main() {
    gl_Position = projection * view * model * vec4(aPosition, 1.0);
    FragPos = vec3(model * vec4(aPosition, 1.0));
    ModelViewPos = vec3(view  * model * vec4(aPosition, 1.0));

    TexCoord = aTexCoord;
    Normal = normalize(mat3(transpose(inverse(view * model))) * aNormal);

    for (int i = 0; i < x_numPointLights; i++ ) {
        LightPositions[i] = vec3(view * vec4(x_lightPositions[i], 1));
    }
}
