import React from 'react';
import {mount, ReactWrapper, shallow} from 'enzyme';
import Table, {IColumn} from './Table';
import {mockSelectionHook} from 'src/test/mockHooks';
import TableHead from '@material-ui/core/TableHead';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';
import TableRow from '@material-ui/core/TableRow';

class TestRow {
    constructor(
        readonly id: string,
        readonly c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

function expectCellText(row: ReactWrapper, ...text: string[]) {
    expect(row.find(TableCell).map((cell) => cell.text())).toEqual(text);
}

function expectCellClasses(row: ReactWrapper, ...classes: string[]) {
    expect(row.find(TableCell).map((cell) => cell.prop('className'))).toEqual(classes);
}

describe('Table', () => {
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    const data = [
        new TestRow('r1', 'a1', 'a2', 'a3'),
        new TestRow('r2', 'b1', 20, 'b3'),
    ];

    it('tracks selected cell', () => {
        const selection = mockSelectionHook();

        const component = shallow(<Table columns={columns} data={data} />);

        const container = component.find('.scroll-container');
        expect(container).toExist();
        expect(container).toHaveProp('onKeyDown', selection.onKeyDown);
        expect(container).toHaveProp('onMouseDown', selection.onMouseDown);
        expect(container).toHaveProp('tabIndex', 0);
    });
    it('displays header row', () => {
        const component = mount(<Table columns={columns} data={data} />);

        const header = component.find(TableHead);

        expectCellText(header, 'first', 'second', 'THIRD');
        expectCellClasses(header, 'col1', 'undefined', '');
    });
    it('populates table body using data', () => {
        const component = mount(<Table columns={columns} data={data} />);

        const body = component.find(TableBody);

        const rows = body.find(TableRow);
        expect(rows).toHaveLength(data.length);
        expect(rows.at(0)).toHaveClassName('selected');
        expectCellText(rows.at(0), 'a1', 'a2', 'a3');
        expectCellClasses(rows.at(0), 'col1 selected', 'string', '');
        expect(rows.at(1)).not.toHaveClassName('selected');
        expectCellText(rows.at(1), 'b1', '20', 'b3');
        expectCellClasses(rows.at(1), 'col1 selected', 'number', '');
    });
});
