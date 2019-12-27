import React from 'react';
import classNames from 'classnames';
import {translate} from '../i18n/localize';
import {IColumn, IRow, ITableProps, getClassName} from './Table';
import ScrollViewport from './ScrollViewport';

export interface IHeaderDetailTableProps<T, S> extends ITableProps<T> {
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
}

type TableType = React.FC<IHeaderDetailTableProps<any, any>>;
const chunk = 50;

const HeaderDetailTable: TableType = <T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, data, subrows, className} = props;
    const [startRow, setStartRow] = React.useState(0);
    const [endRow, setEndRow] = React.useState(0);
    const [offset, setOffset] = React.useState(0);
    const tableRef: React.MutableRefObject<HTMLTableElement> = React.useRef(null);
    const bodyRef: React.MutableRefObject<HTMLTableSectionElement> = React.useRef(null);
    React.useEffect(() => {
        if (data.length > 0 && bodyRef.current.getBoundingClientRect().height === 0) {
            setStartRow(data.length - chunk);
            setEndRow(data.length);
        }
        else if (data.length > 0) {
            const headerHeight = bodyRef.current.getBoundingClientRect().top - tableRef.current.getBoundingClientRect().top;
            const bodyHeight = bodyRef.current.getBoundingClientRect().height;
            const scrollHeight = tableRef.current.parentElement.getBoundingClientRect().height;
            setOffset(scrollHeight - bodyHeight - headerHeight);
        }
    }, [data, tableRef.current, bodyRef.current, startRow, endRow]);
    const onWheel = React.useCallback(({deltaY}: React.WheelEvent) => setOffset((value) => value - deltaY), []);
    return (
        <ScrollViewport scroll={{onWheel, itemCount: data.length, start: data.length / 2, end: data.length / 2 + 20}}>
            <table ref={tableRef} className={classNames('table header-detail', className)} style={{top: offset}}>
                <thead>
                    <tr>
                        {columns.map(({key, className: colClass, colspan, header = translate}) =>
                            <th key={key} className={getClassName(colClass)} colSpan={colspan}>{header(key)}</th>
                        )}
                    </tr>
                    <tr className='detail'>
                        {subColumns.map(({key, className: colClass, colspan, header = translate}) =>
                            <th key={key} className={getClassName(colClass)} colSpan={colspan}>{header(key)}</th>
                        )}
                    </tr>
                </thead>
                <tbody ref={bodyRef}>
                    {data.slice(startRow, endRow).map((row, index) =>
                        <React.Fragment key={row.id}>
                            <tr className={classNames({even: index % 2})}>
                                {columns.map(({key, className: colClass, render, colspan}) =>
                                    <td key={key} className={getClassName(colClass, row)} colSpan={colspan}>{render(row)}</td>
                                )}
                            </tr>
                            {subrows(row).map(subrow =>
                                <tr key={subrow.id} className={classNames('detail', {even: index % 2})}>
                                    {subColumns.map(({key, className: colClass, render, colspan}) =>
                                        <td key={key} className={getClassName(colClass, subrow)} colSpan={colspan}>{render(subrow)}</td>
                                    )}
                                </tr>
                            )}
                        </React.Fragment>
                    )}
                </tbody>
            </table>
        </ScrollViewport>
    );
};

export default HeaderDetailTable;