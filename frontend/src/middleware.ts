import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const publicRoutes = ['/login', '/register', '/verify', '/create-password', '/forgot-password', '/reset-password'];

export function middleware(request: NextRequest) {
  const { nextUrl, cookies } = request;
  const token = cookies.get('auth_token')?.value;

  const isPublicRoute = publicRoutes.some((route) => nextUrl.pathname.startsWith(route));
  const isStaticResource = nextUrl.pathname.startsWith('/_next') || 
                           nextUrl.pathname.startsWith('/favicon.ico');

  // 1. If no token and accessing a protected route -> redirect to login
  if (!token && !isPublicRoute && !isStaticResource) {
    const loginUrl = new URL('/login', request.url);
    return NextResponse.redirect(loginUrl);
  }

  // 2. If token exists and accessing login/register -> redirect to homepage
  if (token && (nextUrl.pathname === '/login' || nextUrl.pathname === '/register')) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
