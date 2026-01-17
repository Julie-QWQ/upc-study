-- 回滚:恢复为英文名称
UPDATE material_categories SET name = 'Courseware', description = 'Lecture slides and courseware' WHERE code = 'courseware';
UPDATE material_categories SET name = 'Textbook', description = 'Textbooks and course books' WHERE code = 'textbook';
UPDATE material_categories SET name = 'Reference', description = 'Reference books and materials' WHERE code = 'reference';
UPDATE material_categories SET name = 'Exam Paper', description = 'Past exam papers' WHERE code = 'exam_paper';
UPDATE material_categories SET name = 'Exercise', description = 'Practice exercises' WHERE code = 'exercise';
UPDATE material_categories SET name = 'Experiment', description = 'Lab guides and experiment reports' WHERE code = 'experiment';
UPDATE material_categories SET name = 'Note', description = 'Class notes and summaries' WHERE code = 'note';
UPDATE material_categories SET name = 'Thesis', description = 'Academic papers and theses' WHERE code = 'thesis';
UPDATE material_categories SET name = 'Other', description = 'Other materials' WHERE code = 'other';
