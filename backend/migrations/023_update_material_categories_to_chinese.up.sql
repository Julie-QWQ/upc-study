-- 更新 material_categories 表中的名称为中文
UPDATE material_categories SET name = '课件资料', description = '教学课件、PPT、讲义等' WHERE code = 'courseware';
UPDATE material_categories SET name = '教材教辅', description = '教材、参考书、教辅资料等' WHERE code = 'textbook';
UPDATE material_categories SET name = '参考文献', description = '参考书目、学术资料等' WHERE code = 'reference';
UPDATE material_categories SET name = '试卷习题', description = '历年试卷、练习题、习题集等' WHERE code = 'exam_paper';
UPDATE material_categories SET name = '练习资料', description = '课后练习、作业题等' WHERE code = 'exercise';
UPDATE material_categories SET name = '实验指导', description = '实验手册、实验报告、指导书等' WHERE code = 'experiment';
UPDATE material_categories SET name = '课堂笔记', description = '课堂笔记、复习总结等' WHERE code = 'note';
UPDATE material_categories SET name = '学术论文', description = '论文、研究报告、学术文章等' WHERE code = 'thesis';
UPDATE material_categories SET name = '其他资料', description = '其他类型的资料' WHERE code = 'other';
