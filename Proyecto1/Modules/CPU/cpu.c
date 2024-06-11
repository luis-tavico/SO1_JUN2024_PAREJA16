// THIS_MODULE, MODULE_VERSION, ...
#include <linux/module.h>
// module_{init,exit}
#include <linux/init.h> 
#include <linux/proc_fs.h>
// for_each_process()
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/fs.h>
#include <linux/sched.h>
// get_mm_rss()
#include <linux/mm.h>
// struct cred, kuid_t
#include <linux/cred.h>
// from_kuid()
#include <linux/uidgid.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Grupo16");
MODULE_DESCRIPTION("Modulo de CPU, Laboratorio Sistemas Operativos 1");

struct task_struct *task;       // sched.h para tareas/procesos
struct task_struct *task_child; // index de tareas secundarias
struct list_head *list;         // lista de cada tareas

static int getPercentageCpu(void)
{
    struct file *file_proc;
    char lectura[256];

    int usuario, nice, system, idle, iowait, irq, softirq, steal, guest, guest_nice;
    int total;
    int porcentaje;

    file_proc = filp_open("/proc/stat", O_RDONLY, 0);
    if (IS_ERR(file_proc))
    {
        printk(KERN_ALERT "Error al abrir el file_proc");
        return -1;
    }

    memset(lectura, 0, 256);
    kernel_read(file_proc, lectura, sizeof(lectura), &file_proc->f_pos);

    sscanf(lectura, "cpu %d %d %d %d %d %d %d %d %d %d", &usuario, &nice, &system, &idle, &iowait, &irq, &softirq, &steal, &guest, &guest_nice);

    total = usuario + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice;

    porcentaje = (total - idle) * 100 / total;
    filp_close(file_proc, NULL);

    return porcentaje;
}

static int escribir_a_proc(struct seq_file *file_proc, void *v)
{
    int porcentajecpu = getPercentageCpu();
    int running = 0;
    int sleeping = 0;
    int zombie = 0;
    int stopped = 0;
    unsigned long rss;
    unsigned long total_ram_pages;
    unsigned long total_usage = 0;
    int b = 0;
    int a = 0;
    int porcentaje;

    if (porcentajecpu == -1)
    {
        seq_printf(file_proc, "Error al leer el archivo");
        return 0;
    }

    total_ram_pages = totalram_pages();
    if (!total_ram_pages) {
        pr_err("Memoria no disponible\n");
        return -EINVAL;
    }

    #ifndef CONFIG_MMU
        pr_err("No MMU, no se puede calcular el RSS.\n");
        return -EINVAL;
    #endif

    for_each_process(task) {
        unsigned long cpu_time = jiffies_to_msecs(task->utime + task->stime);
        total_usage += cpu_time;
    }

    seq_printf(file_proc, "{\n\"cpu_percentage\":%d,\n", porcentajecpu);
    seq_printf(file_proc, "\"processes\":[\n");

    for_each_process(task)
    {
        if (task->mm)
        {
            rss = get_mm_rss(task->mm) << PAGE_SHIFT;
        }
        else
        {
            rss = 0;
        }
        if (b == 0)
        {
            seq_printf(file_proc, "{");
            b = 1;
        }
        else
        {
            seq_printf(file_proc, ",{");
        }
        seq_printf(file_proc, "\"pid\":%d,\n", task->pid);
        seq_printf(file_proc, "\"name\":\"%s\",\n", task->comm);
        seq_printf(file_proc, "\"user\": %u,\n", from_kuid(&init_user_ns, task->cred->uid));
        seq_printf(file_proc, "\"state\":%u,\n", task->__state);
        porcentaje = (rss * 100) / total_ram_pages;
        seq_printf(file_proc, "\"ram\":%d,\n", porcentaje);

        seq_printf(file_proc, "\"child\":[\n");
        a = 0;
        list_for_each(list, &(task->children))
        {
            task_child = list_entry(list, struct task_struct, sibling);
            if (a != 0)
            {
                seq_printf(file_proc, ",{");
            }
            else
            {
                seq_printf(file_proc, "{");
                a = 1;
            }
            seq_printf(file_proc, "\"pid\":%d,\n", task_child->pid);
            seq_printf(file_proc, "\"name\":\"%s\",\n", task_child->comm);
            seq_printf(file_proc, "\"state\":%u,\n", task_child->__state);
            seq_printf(file_proc, "\"pidPadre\":%d\n", task->pid);
            seq_printf(file_proc, "}\n");
        }
        seq_printf(file_proc, "\n]");
        if (task->__state == 0)
        {
            running += 1;
        }
        else if (task->__state == 1)
        {
            sleeping += 1;
        }
        else if (task->__state == 4)
        {
            zombie += 1;
        }
        else
        {
            stopped += 1;
        }
        seq_printf(file_proc, "}\n");
    }

    seq_printf(file_proc, "],\n");
    seq_printf(file_proc, "\"running\":%d,\n", running);
    seq_printf(file_proc, "\"sleeping\":%d,\n", sleeping);
    seq_printf(file_proc, "\"zombie\":%d,\n", zombie);
    seq_printf(file_proc, "\"stopped\":%d,\n", stopped);
    seq_printf(file_proc, "\"total\":%d\n", running + sleeping + zombie + stopped);
    seq_printf(file_proc, "}\n");

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
    proc_create("cpu_so1_1s2024", 0, NULL, &archivo_operaciones);
    printk(KERN_INFO "Modulo CPU creado.\n");
    return 0;
}

static void __exit modulo_cleanup(void)
{
    remove_proc_entry("cpu_so1_1s2024", NULL);
    printk(KERN_INFO "Modulo CPU eliminado.\n");
}

module_init(modulo_init);
module_exit(modulo_cleanup);
