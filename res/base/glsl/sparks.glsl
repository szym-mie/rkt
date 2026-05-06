//+vert
#version 330 core

layout (location = 0) in vec4 a_Pos;

out float v_Decay;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;

void main() {
    gl_Position = u_PMatrix * u_VMatrix * vec4(a_Pos.xyz, 1.0);
    gl_PointSize = 16.0 / (gl_Position.z * 0.05 + 0.1);
    v_Decay = a_Pos.w;
}
//+frag
#version 330 core
precision mediump float;

in float v_Decay;

out vec4 f_Color;

void main() {
    float ttl = 1.0 - v_Decay;
    float hot = ttl * ttl;
    vec3 color = mix(vec3(1.0, 0.2, 0.0), vec3(1.0, 0.8, 0.4), hot);
    f_Color = vec4(color.rgb, ttl);
}