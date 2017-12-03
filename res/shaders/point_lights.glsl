
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

// if(isOutsideSpotLight(view, ligthDirection, spotLight.direction, spotLight.cutoff)) {
//           fragColor = vec4(0);
//           return;
//       }
bool isOutsideSpotLight(mat4 view, vec3 lightDir, vec3 direction, float cutOff) {
    vec3 viewDirection = (view * vec4(direction, 0)).xyz;
    float theta = dot(lightDir, normalize(-viewDirection));
    return theta < cutOff;
}


