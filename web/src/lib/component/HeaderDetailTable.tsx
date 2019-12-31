import React from 'react';
import {observer} from 'mobx-react-lite';
import classNames from 'classnames';
import {translate} from '../i18n/localize';
import {IColumn, IRow, getClassName} from './Table';
import IMixedRowTableModel from '../model/IMixedRowTableModel';
import ScrollViewport, {IScrollableProps} from './ScrollViewport';

export interface IHeaderDetailTableProps<T, S> {
    className?: string;
    columns: IColumn<T>[];
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
    model: IMixedRowTableModel<T>;
}

type TableType = React.FC<IHeaderDetailTableProps<any, any>>;
const defaultRowHeight = 22;

const getHeight = (tableRef: HTMLTableElement, selector: string, defaultHeight: number = 0) => {
    return tableRef ? tableRef.querySelector(selector).clientHeight : defaultHeight;
};

const HeaderDetailTable: TableType = observer(<T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, model, subrows, className} = props;
    const [startRow, setStartRow] = React.useState(0);
    const tableRef: React.MutableRefObject<HTMLTableElement> = React.useRef(null);
    const rowHeight = getHeight(tableRef.current, 'tbody tr.prototype', defaultRowHeight);
    const headerHeight = getHeight(tableRef.current, 'thead');
    // React.useEffect(() => {
    //     if (tableRef.current && model.groups.length > 0) gotoBottom();
    // }, [model, tableRef.current, rowHeight]);
    const startGroup = model.getGroupIndex(startRow);
    const leadingHeight = model.precedingRows[startGroup] * rowHeight;
    const height = model.rowCount * rowHeight + headerHeight;
    const onScroll = React.useCallback(({currentTarget}: React.UIEvent<HTMLElement>) => {
        const {scrollTop} = currentTarget;
        setStartRow(Math.floor(scrollTop / rowHeight));
    }, [rowHeight]);
    return (
        <ScrollViewport onScroll={onScroll}>
            {({scrollHeight}: IScrollableProps) => {
                const endGroup = model.getGroupIndex(startRow + Math.ceil(scrollHeight / rowHeight));
                const trailingHeight = model.getRowsAfter(endGroup) * rowHeight;
                return (
                    <table ref={tableRef} className={classNames('table header-detail', className)} style={{height}}>
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
                        <tbody>
                            <tr className='prototype'><td>0</td></tr>
                            {leadingHeight > 0 ? <tr style={{height: leadingHeight}}><td /></tr> : null}
                            {model.groups.slice(startGroup, endGroup + 1).map((row, index) =>
                                <React.Fragment key={row.id}>
                                    <tr className={classNames({even: (startGroup + index) % 2})}>
                                        {columns.map(({key, className: colClass, render, colspan}) =>
                                            <td key={key} className={getClassName(colClass, row)} colSpan={colspan}>{render(row)}</td>
                                        )}
                                    </tr>
                                    {subrows(row).map(subrow =>
                                        <tr key={subrow.id} className={classNames('detail', {even: (startGroup + index) % 2})}>
                                            {subColumns.map(({key, className: colClass, render, colspan}) =>
                                                <td key={key} className={getClassName(colClass, subrow)} colSpan={colspan}>{render(subrow)}</td>
                                            )}
                                        </tr>
                                    )}
                                </React.Fragment>
                            )}
                            {trailingHeight > 0 ? <tr style={{height: trailingHeight}}><td /></tr> : null}
                        </tbody>
                    </table>
                );
            }}
        </ScrollViewport>
    );
});

export default HeaderDetailTable;