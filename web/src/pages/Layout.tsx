import { Outlet, Link } from 'react-router-dom';
import React, { useState } from 'react';
import {
    BarChartOutlined,
    TableOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Breadcrumb, Layout, Menu, theme } from 'antd';

const { Header, Content, Sider } = Layout;

const headItems: MenuProps['items'] = ['1', '2', '3'].map((key) => ({
    key,
    label: `nav ${key}`,
}));

const sideItems: MenuProps['items'] = [
    {
        key: '/',
        icon: React.createElement(BarChartOutlined),
        label: <Link to="/">Home</Link>,
    },
    {
        key: '/access',
        icon: React.createElement(TableOutlined),
        label: <Link to="/access">Access</Link>,
    },
];

function AppLayout() {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    return (
        <Layout>
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <div className="demo-logo" />
                <Menu
                    theme="dark"
                    mode="horizontal"
                    defaultSelectedKeys={['2']}
                    items={headItems}
                    style={{ flex: 1, minWidth: 0 }}
                />
            </Header>
            <Layout style={{ minHeight: '100vh', background: colorBgContainer }} >
                <Sider theme="light" collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                    <Menu defaultSelectedKeys={['1']} mode="inline" items={sideItems}>
                    </Menu>
                </Sider>
                <Layout style={{ padding: '0 24px 24px' }}>
                    <Breadcrumb style={{ margin: '16px 0' }} items={[{ title: 'User' }, { title: 'Bill' }]} />
                    <Content
                        style={{
                            padding: 24,
                            margin: 0,
                            minHeight: 280,
                            background: colorBgContainer,
                            borderRadius: borderRadiusLG,
                        }}
                    >
                        <Outlet />
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    )
}

export default AppLayout