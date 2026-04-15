//+vert
#version 120

in vec3 a_Pos;
// layout (location = 0) in vec3 a_Pos;

out vec3 v_Pos;

uniform mat4 u_VPIMatrix;

void main() {
    gl_Position = vec4(a_Pos.xy, 1.0, 1.0);
    v_Pos = a_Pos;
}
//+frag
#version 120
precision mediump float;

out vec4 f_Color;

in vec3 v_Pos;

uniform mat4 u_VPIMatrix;
uniform samplerCube u_DiffTexture;

void main() {
    vec4 backVec = u_VPIMatrix * vec4(v_Pos.xyz, 1.0);
    f_Color = textureCube(u_DiffTexture, backVec.xyz / backVec.w);
}