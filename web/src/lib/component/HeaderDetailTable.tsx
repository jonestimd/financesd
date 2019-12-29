import React from 'react';
import classNames from 'classnames';
import {translate} from '../i18n/localize';
import {IColumn, IRow, getClassName} from './Table';
import VirtualScroll from './VirtualScroll';
import IMixedRowTableModel from '../model/IMixedRowTableModel';

export interface IHeaderDetailTableProps<T, S> {
    className?: string;
    columns: IColumn<T>[];
    subColumns: IColumn<S>[];
    subrows: (row: T) => S[];
    model: IMixedRowTableModel<T>;
}

type TableType = React.FC<IHeaderDetailTableProps<any, any>>;
const defaultRowHeight = 22;

const useHeight = (tableRef: HTMLTableElement, selector: string, defaultHeight: number = 0) => {
    return React.useMemo(() => tableRef ? tableRef.querySelector(selector).clientHeight : defaultHeight, [tableRef]);
};

const useRedraw = () => {
    const [redraw, setRedraw] = React.useState(false);
    React.useEffect(() => setRedraw(false), [redraw]);
    return React.useCallback(() => setRedraw(true), []);
};

const HeaderDetailTable: TableType = <T extends IRow, S extends IRow>(props: IHeaderDetailTableProps<T, S>) => {
    const {columns, subColumns, model, subrows, className} = props;
    const [startRow, setStartRow] = React.useState(0);
    const [offset, setOffset] = React.useState(0);
    const onResize = useRedraw();
    const tableRef: React.MutableRefObject<HTMLTableElement> = React.useRef(null);
    const rowHeight = useHeight(tableRef.current, 'tbody tr.prototype', defaultRowHeight);
    const headerHeight = useHeight(tableRef.current, 'thead');
    const scrollHeight = tableRef.current ? tableRef.current.parentElement.clientHeight - headerHeight : 0;
    React.useEffect(() => {
        if (tableRef.current && model.groups.length > 0) {
            const visibleRows = Math.ceil(scrollHeight / rowHeight);
            const start = Math.max(0, model.rowCount - visibleRows);
            const initOffset = Math.min(0, scrollHeight - (model.rowCount - start) * rowHeight);
            setStartRow(start);
            setOffset(initOffset);
        }
    }, [model, tableRef.current, rowHeight]);
    const [startGroup, precedingRows] = model.getGroupIndex(startRow);
    const endRow = startRow + Math.ceil(scrollHeight / rowHeight);
    const [endGroup] = model.getGroupIndex(endRow);
    const top = offset - (startRow - precedingRows) * rowHeight;
    const onScroll = React.useCallback((deltaY: number) => {
        let newOffset = offset - deltaY, start = startRow;
        if (newOffset > 0 && start > 0) {
            const deltaRows = Math.ceil(newOffset / rowHeight);
            start = start - deltaRows;
        }
        else if (newOffset < -rowHeight) {
            const deltaRows = Math.floor(-newOffset / rowHeight);
            start = start + deltaRows;
        }
        if (start >= 0 && start < model.rowCount) {
            newOffset += rowHeight * (start - startRow);
            setStartRow(start);
            setOffset(newOffset);
        }
    }, [offset, startRow, rowHeight]);
    return (
        <VirtualScroll onScroll={onScroll} onResize={onResize} itemCount={model.rowCount} start={startRow} end={endRow}>
            <table ref={tableRef} className={classNames('table header-detail', className)} style={{top}}>
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
                </tbody>
            </table>
        </VirtualScroll>
    );
};

export default HeaderDetailTable;