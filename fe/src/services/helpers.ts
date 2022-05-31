import { Tag } from 'models/Tag';

export const getArrayOfObjects = (data: Tag[]) => {
    return data.map((tag: Tag) => {
        return {
            [tag.key]: tag.value,
        };
    });
};

export const getResourcesRequestData = (data: Tag[]) => {
    const tags: {
        [key: string]: string;
    } = {};

    data.forEach((tag: Tag) => {
        tags[tag.key] = tag.value;
    });

    const filter = tags;

    return { filter };
};
