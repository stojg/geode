#include "include.glsl"

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec2 TexCoord;
out vec3 Normal;

void setOutputs() {
    gl_Position = projection * view * model * vec4(aPosition, 1.0);
    Normal = aNormal;
    TexCoord = aTexCoord;
}
