import ChangeModel from './ChangeModel';

describe('ChangeModel', () => {
    describe('get', () => {
        it('returns the original value if no change', () => {
            const model = new ChangeModel({value: 'original'});

            expect(model.get('value')).toEqual('original');
        });
        it('returns change value', () => {
            const model = new ChangeModel({value: 'original'});

            model.set('value', 'new data');

            expect(model.get('value')).toEqual('new data');
        });
        it('returns undefined for cleared value', () => {
            const model = new ChangeModel({value: 'original'});

            model.set('value', undefined);

            expect(model.get('value')).toBeUndefined();
        });
    });
    describe('set', () => {
        it('saves change', () => {
            const model = new ChangeModel({value: 'original'});

            model.set('value', 'new data');

            expect(model.changes).toEqual({value: 'new data'});
            expect(model.isChanged).toBe(true);
        });
        it('clears value', () => {
            const model = new ChangeModel({value: 'original'});

            model.set('value', undefined);

            expect(model.changes).toEqual({value: undefined});
            expect(model.isChanged).toBe(true);
        });
        it('resets change', () => {
            const model = new ChangeModel({value: 'original'});
            model.set('value', 'new data');

            model.set('value', 'original');

            expect(model.changes).toEqual({});
            expect(model.isChanged).toBe(false);
        });
    });
    describe('undo', () => {
        it('removes change', () => {
            const model = new ChangeModel({value: 'original'});
            model.set('value', undefined);

            model.undo('value');

            expect(model.get('value')).toEqual('original');
            expect(model.changes).toEqual({});
            expect(model.isChanged).toBe(false);
        });
    });
    describe('commit', () => {
        it('updates original values', () => {
            const model = new ChangeModel({value: 'original'});
            model.set('value', 'new data');

            model.commit();

            expect(model.isChanged).toBe(false);
            expect(model.get('value')).toEqual('new data');
        });
    });
    describe('revert', () => {
        it('clears changes', () => {
            const model = new ChangeModel({value: 'original'});
            model.set('value', undefined);

            model.revert();

            expect(model.get('value')).toEqual('original');
            expect(model.changes).toEqual({});
            expect(model.isChanged).toBe(false);
        });
    });
});
