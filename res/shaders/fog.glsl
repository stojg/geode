const vec3 fogColor = vec3(0.9, 0.9,1);
const float FogDensity = 0.003;

vec3 fogCalc(vec3 Lo, vec3 vp) {
    float dist = 0;
    float fogFactor = 0;
    dist = length(vp);
    fogFactor = 1.0 /exp(dist * FogDensity);
    fogFactor = clamp( fogFactor, 0.0, 1.0 );
    return mix(fogColor, Lo, fogFactor);
}
