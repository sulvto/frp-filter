import { useEffect, useRef, useState } from 'react'
import type { TableColumnsType, TableProps, GetProp, InputRef, TableColumnType } from 'antd';
import { Button, Input, Space, Table, Skeleton } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import type { SorterResult } from 'antd/es/table/interface';
import type { FilterDropdownProps } from 'antd/es/table/interface';
import Highlighter from 'react-highlight-words';

function Access() {
    interface DataType {
        no: number;
        ip: string;
        time: string;
        location: string;
        count: number;
    }
    type DataIndex = keyof DataType;
    type TablePaginationConfig = Exclude<GetProp<TableProps, 'pagination'>, boolean>;

    interface TableParams {
        pagination?: TablePaginationConfig;
        sortField?: SorterResult<any>['field'];
        sortOrder?: SorterResult<any>['order'];
        filters?: Parameters<GetProp<TableProps, 'onChange'>>[1];
    }

    const [data, setData] = useState<DataType[]>();
    const [searchText, setSearchText] = useState('');
    const [searchedColumn, setSearchedColumn] = useState('');
    const searchInput = useRef<InputRef>(null);
    const [loading, setLoading] = useState(false);
    const [initializing, setInitializing] = useState(true);
    const [tableParams, setTableParams] = useState<TableParams>({
        pagination: {
            current: 1,
            pageSize: 10,
        },
    });

    const handleSearch = (
        selectedKeys: string[],
        confirm: FilterDropdownProps['confirm'],
        dataIndex: DataIndex,
    ) => {
        confirm();
        setSearchText(selectedKeys[0]);
        setSearchedColumn(dataIndex);
    };

    const handleReset = (clearFilters: () => void) => {
        clearFilters();
        setSearchText('');
    };

    const getColumnSearchProps = (dataIndex: DataIndex): TableColumnType<DataType> => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters, close }) => (
            <div style={{ padding: 8 }} onKeyDown={(e) => e.stopPropagation()}>
                <Input
                    ref={searchInput}
                    placeholder={`Search ${dataIndex}`}
                    value={selectedKeys[0]}
                    onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
                    onPressEnter={() => handleSearch(selectedKeys as string[], confirm, dataIndex)}
                    style={{ marginBottom: 8, display: 'block' }}
                />
                <Space>
                    <Button
                        type="primary"
                        onClick={() => handleSearch(selectedKeys as string[], confirm, dataIndex)}
                        icon={<SearchOutlined />}
                        size="small"
                        style={{ width: 90 }}
                    >
                        Search
                    </Button>
                    <Button
                        onClick={() => clearFilters && handleReset(clearFilters)}
                        size="small"
                        style={{ width: 90 }}
                    >
                        Reset
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        onClick={() => {
                            confirm({ closeDropdown: false });
                            setSearchText((selectedKeys as string[])[0]);
                            setSearchedColumn(dataIndex);
                        }}
                    >
                        Filter
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        onClick={() => {
                            close();
                        }}
                    >
                        close
                    </Button>
                </Space>
            </div>
        ),
        filterIcon: (filtered: boolean) => (
            <SearchOutlined style={{ color: filtered ? '#1677ff' : undefined }} />
        ),
        onFilter: (value, record) => record[dataIndex]
            .toString()
            .toLowerCase()
            .includes((value as string).toLowerCase()),
        onFilterDropdownOpenChange: (visible) => {
            if (visible) {
                setTimeout(() => searchInput.current?.select(), 100);
            }
        },
        render: (text) =>
            searchedColumn === dataIndex ? (
                <Highlighter
                    highlightStyle={{ backgroundColor: '#ffc069', padding: 0 }}
                    searchWords={[searchText]}
                    autoEscape
                    textToHighlight={text ? text.toString() : ''}
                />
            ) : (
                text
            ),
    });

    const toLocation = (info: any) => {
        if (info) {
            let area_arr: string[] = [info.countryCN, info.provinceCN, info.cityCN];
            let area: string = area_arr.filter((item, index) => area_arr.indexOf(item) === index).filter(item => item !== '*').join('-');
            return `${area} ${info.ispCN}`.trim();
        } else {
            return "";
        }
    }

    const fetchIpInfo = (record: DataType) => {
        fetch(`http://211.149.239.251:7777/ip/location?ip=${record.ip}`)
            .then((res) => res.json())
            .then((info) => {
                let location = toLocation(info)
                setData((data) =>
                    data ? data.map(item =>
                        item.no === record.no ? { ...item, location: location } : item
                    ) : []
                );
            });
    };

    const fetchData = () => {
        setLoading(true);
        fetch(`http://211.149.239.251:7777/access`)
            .then((res) => res.json())
            .then((data) => {
                for (var i = 0; i < data.length; i++) {
                    data[i].no = i;
                    data[i].location = toLocation(data[i].info);
                }
                setData(data);
                setLoading(false);
                setInitializing(false);
                setTableParams({
                    ...tableParams,
                    pagination: {
                        ...tableParams.pagination,
                        total: data.length,
                    },
                });
            });
    };

    useEffect(fetchData, [
        tableParams.pagination?.current,
        tableParams.pagination?.pageSize,
        tableParams?.sortOrder,
        tableParams?.sortField,
        JSON.stringify(tableParams.filters),
    ]);

    const handleTableChange: TableProps<DataType>['onChange'] = (pagination, filters, sorter) => {
        setTableParams({
            pagination,
            filters,
            sortOrder: Array.isArray(sorter) ? undefined : sorter.order,
            sortField: Array.isArray(sorter) ? undefined : sorter.field,
        });

        // `dataSource` is useless since `pageSize` changed
        if (pagination.pageSize !== tableParams.pagination?.pageSize) {
            setData([]);
        }
    };

    const columns: TableColumnsType<DataType> = [
        {
            title: 'IP',
            dataIndex: 'ip',
            key: 'ip',
            onFilter: (value, record) => record.ip.indexOf(value as string) === 0,
            ...getColumnSearchProps('ip'),
        },
        {
            title: '访问时间',
            dataIndex: 'time',
            key: 'time',
            onFilter: (value, record) => record.time.indexOf(value as string) === 0,
            sorter: (a, b) => a.time.localeCompare(b.time),
            sortDirections: ['ascend', 'descend'],
        },
        {
            title: '归属地',
            dataIndex: 'location',
            key: 'location',
        },
        {
            title: '访问次数',
            dataIndex: 'count',
            key: 'count',
            sorter: (a, b) => a.count - b.count,
            sortDirections: ['ascend', 'descend'],
        },
        {
            title: '操作',
            width: 250,
            fixed: 'right',
            render: (_, record) => (
                <Space>
                    <a onClick={() => fetchIpInfo(record)}>查询归属地</a>
                    <a>屏蔽</a>
                    <a>取消屏蔽</a>
                    <a>删除</a>
                </Space>
            ),
        },
    ];

    return (
        <Skeleton active loading={initializing} >
            <Table<DataType>
                dataSource={data}
                columns={columns}
                rowKey={(record) => record.no}
                pagination={tableParams.pagination}
                loading={loading}
                onChange={handleTableChange}
            />
        </Skeleton>
    )
}

export default Access
