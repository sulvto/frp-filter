import { Card, Flex } from 'antd';
import { Column } from '@ant-design/plots';
import { Line } from '@ant-design/charts';

const data = [
    { year: '1991', value: 3 },
    { year: '1992', value: 4 },
    { year: '1993', value: 3.5 },
    { year: '1994', value: 5 },
    { year: '1995', value: 4.9 },
    { year: '1996', value: 6 },
    { year: '1997', value: 7 },
    { year: '1998', value: 9 },
    { year: '1999', value: 13 },
];

const props = {
    data,
    xField: 'year',
    yField: 'value',
};


const config = {
    data: {
        type: 'fetch',
        value: 'https://render.alipay.com/p/yuyan/180020010001215413/antd-charts/column-column.json',
    },
    xField: 'letter',
    yField: 'frequency',
    label: {
        text: (d: any) => `${(d.frequency * 100).toFixed(1)}%`,
        textBaseline: 'bottom',
    },
    axis: {
        y: {
            labelFormatter: '.0%',
        },
    }
};

export function Overview() {
    return <>
        <Flex gap="middle" vertical>
            <Flex gap="small" vertical={false} >
                <Card style={{
                    width: '50%',
                    height: '100%',
                }}>
                    <Line {...props} />
                </Card>
                <Card style={{
                    width: '50%',
                    height: '100%',
                }}>
                    <Column {...config} />
                </Card>
            </Flex>
        </Flex>

    </>
        ;
}