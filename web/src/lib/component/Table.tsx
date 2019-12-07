import React from 'react';

export interface IColumn<T> {
    name: string;
    className?: string;
    getter: (row: T) => React.ReactNode;
}

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
}

const Table: React.FC<ITableProps<any>> = <T extends {id: string}>({columns, data}: ITableProps<T>) => {
    return (
        <table className='table'>
            <thead>
                <tr>
                    {columns.map((column, i) =>
                        <th key={column.name} className={column.className}>{column.name}</th>
                    )}
                </tr>
            </thead>
            <tbody>
                {data.map(row => (
                    <tr key={row.id}>
                        {columns.map(column =>
                            <td key={column.name} className={column.className}>{column.getter(row)}</td>
                        )}
                    </tr>
                ))}
            </tbody>
        </table>
    );
};

export default Table;