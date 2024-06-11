#include <linux/module.h>
//Header necesario porque se usara proc_fs
#include <linux/proc_fs.h>
#include <linux/sysinfo.h>
/* Header para usar la lib seq_file y manejar el archivo en /proc*/
#include <linux/seq_file.h>
// para get_mm_rss
#include <linux/mm.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Grupo16");
MODULE_DESCRIPTION("Modulo de RAM, Laboratorio Sistemas Operativos 1");

struct sysinfo inf;

static int escribir_a_proc(struct seq_file *file_proc, void *v)
{
    unsigned long total, used, notused;
    unsigned long perc;
    si_meminfo(&inf);

    total = inf.totalram * inf.mem_unit;
    used = inf.freeram * inf.mem_unit + inf.bufferram * inf.mem_unit + inf.sharedram * inf.mem_unit;
    perc = (used * 100) / total;
    notused = total - used;
    seq_printf(file_proc, "{\"ram_total\":%lu, \"ram_used\":%lu, \"ram_percentage\":%lu, \"ram_free\":%lu }", total, used, perc, notused);
    return 0;
}

static int abrir_aproc(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_a_proc, NULL);
}
static struct proc_ops archivo_operaciones = {
    .proc_open = abrir_aproc,
    .proc_read = seq_read
};

static int __init modulo_init(void)
{
    proc_create("ram_so1_jun2024", 0, NULL, &archivo_operaciones);
    printk(KERN_INFO "Laboratorio Sistemas Operativos 1\n");
    return 0;
}

static void __exit modulo_cleanup(void)
{
    remove_proc_entry("ram_so1_jun2024", NULL);
    printk(KERN_INFO "Laboratorio Sistemas Operativos 1\n");
     
}

module_init(modulo_init);
module_exit(modulo_cleanup);