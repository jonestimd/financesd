import React from 'react';
import {shallow} from 'enzyme';
import MessageStore from 'src/lib/store/MessageStore';
import GroupStore from 'src/lib/store/GroupStore';
import Group from './Group';
import {GroupModel} from 'src/lib/model/GroupModel';

describe('Group', () => {
    const groupStore = new GroupStore(new MessageStore());

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({groupStore});
    });
    it('returns null if no group', () => {
        expect(shallow(<Group />)).toBeEmptyRender();
    });
    it('displays group name', () => {
        const id = 234;
        const name = 'the group';
        jest.spyOn(groupStore, 'getGroup').mockReturnValue({name} as GroupModel);

        const component = shallow(<Group id={id} />);

        expect(component.find('.group')).toHaveText(name);
        expect(groupStore.getGroup).toBeCalledWith(id);
    });
});
