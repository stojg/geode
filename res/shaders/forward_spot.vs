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

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 aTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform SpotLight spotLight;

// shadow
uniform mat4 lightViewProjection;
out vec4 FragPosLightSpace;

out vec2 TexCoord;
out vec3 FragPos;
out vec3 Normal;
out vec3 LightPos;

void main() {
    gl_Position = projection * view * model * vec4(position, 1.0);
    FragPos = vec3(view  * model * vec4(position, 1.0));

    Normal = mat3(transpose(inverse(view * model))) * normal;
    TexCoord = aTexCoord;

    LightPos = vec3(view * vec4(spotLight.pointLight.position, 1.0));

    FragPosLightSpace = lightViewProjection * model * vec4(position, 1.0);
}
