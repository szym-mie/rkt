//+vert
#version 330 core

layout (location = 0) in vec3 a_Pos;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;

void main() {
    vec4 outPos = u_PMatrix * u_VMatrix * vec4(a_Pos.xyz, 1.0);
    gl_Position = outPos;
    gl_PointSize = 8.0 / (outPos.z * 0.05 + 0.1);
}
//+frag
#version 330 core
precision mediump float;

out vec4 f_Color;

void main() {
    f_Color = vec4(1.0, 0.5, 0.0, 1.0);
}