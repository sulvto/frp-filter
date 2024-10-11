import { RouterProvider, Navigate, createHashRouter } from 'react-router-dom';
import Error403 from './pages/403.tsx'
import Error404 from './pages/404.tsx'
import Access from './pages/Access'
import AppLayout from './pages/Layout'
import Overview from './pages/Overview'
import Blacklist from './pages/Blacklist'

const routes = [
    {
        path: '/',
        // element: <Navigate to='/welcome' replace />,
        element: <AppLayout />,
        children: [
            {
                path: '',
                element: <Overview />
            },
            {
                path: '/access',
                element: <Access />
            },
            {
                path: '/blacklist',
                element: <Blacklist />
            },
        ]
    },
    {
        path: '/403',
        element: <Error403 />
    },
    {
        path: '/404',
        element: <Error404 />
    },
    {
        path: '*',
        element: <Navigate to={'/404'} replace />
    }
]
const router = createHashRouter(routes)

function App() {
    return (
        <RouterProvider router={router} />
    );
}

export default App
