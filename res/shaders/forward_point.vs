#version 410 core

struct Attenuation
{
    float constant;
    float linear;
    float exponent;
};

struct BaseLight
{
    vec3 color;
};

struct DirectionalLight
{
    BaseLight base;
    vec3 direction;
};

struct PointLight
{
    BaseLight base;
    Attenuation atten;
    vec3 position;
};

struct SpotLight
{
    PointLight pointLight;
    vec3 direction;
    float cutoff;
};

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec2 TexCoord;
out vec3 FragPos;
out vec3 ViewDirection;
out vec3 Normal;
out vec3 LightPos;

// shadow & light
uniform PointLight pointLight;
// uniform mat4 lightMVP; // @todo - implement
// out vec4 FragPosLightSpace;  @todo - implement

void main() {
    gl_Position = projection * view * model * vec4(aPosition, 1.0);
    FragPos = vec3(view  * model * vec4(aPosition, 1.0));
    ViewDirection = -normalize(FragPos);

    Normal = normalize(mat3(transpose(inverse(view * model))) * aNormal);
    TexCoord = aTexCoord;

    LightPos = vec3(view * vec4(pointLight.position, 1.0));
    // FragPosLightSpace = lightMVP * vec4(aPosition, 1.0);
}
