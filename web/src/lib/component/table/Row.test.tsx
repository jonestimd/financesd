import React from 'react';
import {mount, shallow} from 'enzyme';
import Row from './Row';
import {TableCell} from '@material-ui/core';
import {IColumn} from './Column';
import TextCellEditor from './TextCellEditor';
import {Icon} from '@material-ui/core';

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
    const props = {
        columns,
        isEditable: jest.fn(),
    };
    beforeEach(() => {
        // ignore React DOM nesting warnings
        jest.spyOn(console, 'error').mockImplementation(() => {});
    });
    it('populates cells using data', () => {
        const data = new TestRow('r1', 'a1', 'a2', 'a3');

        const component = shallow(<Row {...props} row={data} className='row-class' />);

        const cells = component.find(TableCell);
        expect(component).toHaveClassName('row-class');
        expect(cells.map((c) => c.text())).toEqual(['a1', 'a2', 'a3']);
        expect(cells.map((c) => c.prop('colSpan'))).toEqual([undefined, 2, undefined]);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'string', '']);
    });
    it('highlights selected cell', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = shallow(<Row {...props} row={data} selection={{column: 1}} />);

        const cells = component.find(TableCell);
        expect(cells.map((c) => c.text())).toEqual(['a1', '2', 'a3']);
        expect(cells.map((c) => c.prop('className'))).toEqual(['col1', 'number selected', '']);
    });
    it('displays error icon if validate returns non-empty array', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const validate = (_row: TestRow, key: string) => key === 'second' ? ['invalid'] : [];

        const component = shallow(<Row {...props} row={data} selection={{column: 1}} validate={validate} />);

        expect(component.find(TableCell).at(0).find(Icon)).not.toExist();
        expect(component.find(TableCell).at(1).find(Icon)).toHaveText('error_outline');
        expect(component.find(TableCell).at(2).find(Icon)).not.toExist();
    });
    it('calls onClick', () => {
        const onClick = jest.fn();
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const component = shallow(<Row {...props} row={data} selection={{column: 1}} onClick={onClick} />);

        component.find(TableCell).at(1).simulate('click');

        expect(onClick).toBeCalledWith(1);
    });
    it('disables editing if isEditable returns false', () => {
        props.isEditable.mockReturnValue(false);
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = mount(<Row {...props} row={data} selection={{column: 0}} editCell={0} />);

        expect(component.find(TableCell).at(0).find(TextCellEditor)).not.toExist();
    });
    it('enables editing if no isEditable', () => {
        const {isEditable, ...rest} = props;
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = mount(<Row {...rest} row={data} selection={{column: 0}} editCell={0} />);

        expect(component.find(TableCell).at(0).find(TextCellEditor)).toHaveProp('value', data.c1);
    });
    it('displays cell editor', () => {
        props.isEditable.mockReturnValue(true);
        const data = new TestRow('r1', 'a1', 2, 'a3');

        const component = mount(<Row {...props} row={data} selection={{column: 0}} editCell={0} />);

        expect(component.find(TableCell).at(0).find(TextCellEditor)).toHaveProp('value', data.c1);
    });
    it('updates cell value on commit', () => {
        const data = new TestRow('r1', 'a1', 2, 'a3');
        props.isEditable.mockReturnValue(true);
        const component = mount(<Row {...props} row={data} selection={{column: 0}} editCell={0} />);

        component.find(TextCellEditor).prop('onCommit')('new value');

        expect(editor.setValue).toBeCalledWith(data, 'new value');
    });
    it('calls onCommit on committing cell edit', () => {
        const onCommit = jest.fn();
        props.isEditable.mockReturnValue(true);
        const data = new TestRow('r1', 'a1', 2, 'a3');
        const component = mount(<Row {...props} row={data} selection={{column: 0}} editCell={0} onCommit={onCommit} />);

        component.find(TextCellEditor).prop('onCommit')('new value');

        expect(editor.setValue).toBeCalledWith(data, 'new value');
        expect(onCommit).toBeCalledTimes(1);
    });
});
