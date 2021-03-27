import {newGroup} from 'src/test/groupFactory';
import {GroupModel} from './GroupModel';


describe('GroupModel', () => {
    const group = newGroup();

    describe('constructor', () => {
        it('populates group properties', () => {
            const model = new GroupModel(group);

            expect(model).toEqual(group);
        });
    });
});
