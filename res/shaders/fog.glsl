//const vec3 fogColor = vec3(0.8, 1,0.9);
const vec3 fogColor = vec3(0.416, 0.506, 0.424)*2;
const float FogDensity = 0.001;

vec3 fogCalc(vec3 Lo, vec3 vp) {
    float dist = 0;
    float fogFactor = 0;
    dist = length(vp);
    fogFactor = 1.0 /exp(dist * FogDensity);
    fogFactor = clamp( fogFactor, 0.0, 1.0 );
    return mix(fogColor, Lo, fogFactor);
}
