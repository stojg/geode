
#include "light.glh"

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec2 TexCoord;
out vec3 ModelViewPos;
out vec3 Normal;
out vec3 LightPos;

// shadow & light
uniform mat4 lightMVP;
out vec4 FragPosLightSpace;

void setOutput(vec4 lightPosition) {
gl_Position = projection * view * model * vec4(aPosition, 1.0);
    ModelViewPos = vec3(view  * model * vec4(aPosition, 1.0));

    Normal = normalize(mat3(transpose(inverse(view * model))) * aNormal);
    TexCoord = aTexCoord;

    FragPosLightSpace = lightMVP * vec4(aPosition, 1.0);

    LightPos = vec3(view * lightPosition);
}
