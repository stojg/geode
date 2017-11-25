#version 410 core

out vec4 fragColor;

void main()
{
    float depth = gl_FragCoord.z;
    fragColor = vec4(depth, depth*depth, 0, 0);
    // Adjusting moments (this is sort of bias per pixel) using partial derivative
    float dx = dFdx(depth);
    float dy = dFdx(depth);
    float moment2 = depth * depth + 0.25 * (dx * dx + dy * dy);

    fragColor = vec4(depth, moment2, 0, 0);
}
