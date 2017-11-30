
#include "include.glsl"

in vec2 TexCoord;
in vec3 Normal;

out vec3 FragColor;

uniform sampler2D diffuse;

vec3 calcLight(vec3 color) {
    vec3 result = texture(diffuse, TexCoord).rgb;
    result *= color;
    result *= Normal;
    return result;
}
