import React from 'react';
import classnames from 'classnames';
import {translate} from '../i18n/localize';
import {IColumn, IRow, ITableProps, getClassName} from './Table';

export interface IHeaderDetailTableProps<T, S> extends ITableProps<T> {
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
}

type TableType = React.FC<IHeaderDetailTableProps<any, any>>;

const HeaderDetailTable: TableType = <T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, data, subrows, className} = props;
    return (
        <table className={classnames('table header-detail', className)}>
            <thead>
                <tr>
                    {columns.map(({key, className, colspan, header = translate}) =>
                        <th key={key} className={getClassName(className)} colSpan={colspan}>{header(key)}</th>
                    )}
                </tr>
                <tr className='detail'>
                    {subColumns.map(({key, className, colspan, header = translate}) =>
                        <th key={key} className={getClassName(className)} colSpan={colspan}>{header(key)}</th>
                    )}
                </tr>
            </thead>
            {data.map(row =>
                <tbody key={row.id}>
                    <tr>
                        {columns.map(({key, className, render, colspan}) =>
                            <td key={key} className={getClassName(className, row)} colSpan={colspan}>{render(row)}</td>
                        )}
                    </tr>
                    {subrows(row).map(subrow =>
                        <tr key={subrow.id} className='detail'>
                            {subColumns.map(({key, className, render, colspan}) =>
                                <td key={key} className={getClassName(className, subrow)} colSpan={colspan}>{render(subrow)}</td>
                            )}
                        </tr>
                    )}
                </tbody>
            )}
        </table>
    );
};

export default HeaderDetailTable;