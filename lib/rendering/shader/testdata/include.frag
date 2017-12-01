
#include "include.glsl"

in VS_OUT {
    vec2 TexCoord;
    vec3 Normal;
} vs_in;

out vec3 FragColor;

uniform sampler2D diffuse;

vec3 calcLight(vec3 color) {
    vec3 result = texture(diffuse, vs_in.TexCoord).rgb;
    result *= color;
    result *= vs_in.Normal;
    return result;
}
