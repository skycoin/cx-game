package render

// import (
// 	"github.com/go-gl/gl/v4.1-core/gl"
// 	"github.com/go-gl/mathgl/mgl32"

// 	"github.com/skycoin/cx-game/cxmath"
// )
// var g_display_width = 800;
// var g_display_height = 600;
// var g_colorbits = 32;
// var g_depthbits = 16;
// var g_multisample = 8;
// var g_vsync = true;
// var g_fullscreen = true;
// var g_resize = false;

// // changes display resolution at any time
// bool CApp::ChangeVideoMode()
// {
//  SDL_SysWMinfo info;

//  // get window handle from SDL
//  SDL_VERSION(&info.version);
//  if (SDL_GetWMInfo(&info) == -1) {
//   PrintError("SDL_GetWMInfo #1 failed\n");
//   return false;
//  }

//  // get device context handle
//  HDC tempDC = GetDC( info.window );

//  // create temporary context
//  HGLRC tempRC = wglCreateContext( tempDC );
//  if (tempRC == NULL) {
//   PrintError("wglCreateContext failed\n");
//   return false;
//  }

//  // share resources to temporary context
//  SetLastError(0);
//  if (!wglShareLists(info.hglrc, tempRC)) {
//   PrintError("wglShareLists #1 failed\n");
//   return false;
//  }

//  // set video mode
//  if (!SetVideoMode()) return false;

//  // previously used structure may possibly be invalid, to be sure we get it again
//  SDL_VERSION(&info.version);
//  if (SDL_GetWMInfo(&info) == -1) {
//   PrintError("SDL_GetWMInfo #2 failed\n");
//   return false;
//  }

//  // share resources to new SDL-created context
//  if (!wglShareLists(tempRC, info.hglrc)) {
//   PrintError("wglShareLists #2 failed\n");
//   return false;
//  }

//  // we no longer need our temporary context
//  if (!wglDeleteContext(tempRC)) {
//   PrintError("wglDeleteContext failed\n");
//   return false;
//  }

//  // success
//  return true;
// }
