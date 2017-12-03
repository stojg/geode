
struct Light {
    // common
    vec3 position;
    vec3 color;
    float constant;
    float linear;
    float quadratic;
    float distance;

    // spotlights
    float cutoff;
    vec3 direction;
};
