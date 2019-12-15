import React, {ReactNode} from 'react';
import classnames from 'classnames';
import {translate} from '../i18n/localize';

export interface IColumn<T> {
    key: string;
    className?: string;
    colspan?: number;
    header?: (key: string) => ReactNode;
    render: (row: T) => React.ReactNode;
}

export interface ITableProps<T> {
    columns: IColumn<T>[];
    data: T[];
    className?: string;
}

export interface IRow {
    id: string;
}

const Table: React.FC<ITableProps<any>> = <T extends IRow>({columns, data, className}: ITableProps<T>) => {
    return (
        <table className={classnames('table', className)}>
            <thead>
                <tr>
                    {columns.map(({key, className, header = translate}) => <th key={key} className={className}>{header(key)}</th>)}
                </tr>
            </thead>
            <tbody>
                {data.map(row => (
                    <tr key={row.id}>
                        {columns.map(column =>
                            <td key={column.key} className={column.className}>{column.render(row)}</td>
                        )}
                    </tr>
                ))}
            </tbody>
        </table>
    );
};

export default Table;