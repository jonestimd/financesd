import React from 'react';
import {mount, shallow} from 'enzyme';
import Row from './Row';
import TableCell from '@material-ui/core/TableCell';
import {IColumn} from './Column';
import TextCellEditor from './TextCellEditor';

class TestRow {
    constructor(
        readonly id: string,
        public c1: string,
        readonly c2: string | number,
        readonly c3: string
    ) { }
}

describe('Row', () => {
    const editor = {
        Component: TextCellEditor,
        getValue: (row: TestRow) => row.c1,
        setValue: jest.fn(),
    };
    const columns: IColumn<TestRow>[] = [
        {key: 'first', className: 'col1', render: (row) => row.c1, editor},
        {key: 'second', className: (row) => typeof row?.c2, render: (row) => row.c2, colspan: 2},
        {key: 'third', render: (row) => row.c3, header: (key) => key.toUpperCase()},
    ];
    beforeEach(() => {
        // ignore React DOM nesting warnings
        jest.spyOn(console, 'error').mockImplementation(() => {});
    });
    it('populates cells using data', () => {
        const data = new TestRow('r1', 'a1', 'a2', 'a3');

        const component = shallow(<Row<TestRow> columns={columns} row={data} className='row-class' />);

        const cells = component.find(TableCell);
        expect(component).toHaveClassName('row-class');
        expect(cells.map((c) => c.text())).toEqual(['a1', 'a2', 'a3']);
        expect(cells.map((c) => c.prop('colSpan'))).toEqual([undefined, 2, undefined]);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'string', '']);
    });
    it('highlights selected cell', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = shallow(<Row<TestRow> columns={columns} row={data} selection={{column: 1}} />);

        const cells = component.find(TableCell);
        expect(cells.map((c) => c.text())).toEqual(['a1', '2', 'a3']);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'number selected', '']);
    });
    it('calls onClick', () => {
        const onClick = jest.fn();
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const component = shallow(<Row<TestRow> columns={columns} row={data} selection={{column: 1}} onClick={onClick} />);

        component.find(TableCell).at(1).simulate('click');

        expect(onClick).toBeCalledWith(1);
    });
    it('displays cell editor', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = mount(<Row<TestRow> columns={columns} row={data} selection={{column: 0}} editCell={0} />);

        expect(component.find(TableCell).at(0).find(TextCellEditor)).toHaveProp('value', data.c1);
    });
    it('updates cell value on commit', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const component = mount(<Row<TestRow> columns={columns} row={data} selection={{column: 0}} editCell={0} />);

        component.find(TextCellEditor).prop('onCommit')('new value');

        expect(editor.setValue).toBeCalledWith(data, 'new value');
    });
    it('calls onCommit on committing cell edit', () => {
        const onCommit = jest.fn();
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const component = mount(<Row<TestRow> columns={columns} row={data} selection={{column: 0}} editCell={0} onCommit={onCommit} />);

        component.find(TextCellEditor).prop('onCommit')('new value');

        expect(editor.setValue).toBeCalledWith(data, 'new value');
        expect(onCommit).toBeCalledTimes(1);
    });
});
